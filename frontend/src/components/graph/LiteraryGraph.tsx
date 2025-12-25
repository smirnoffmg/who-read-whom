"use client";

import { OpinionService } from "@/services/opinionService";
import { WorkService } from "@/services/workService";
import { WriterService } from "@/services/writerService";
import type { Work } from "@/types/work";
import type { Writer } from "@/types/writer";
import dynamic from "next/dynamic";
import { useCallback, useEffect, useRef, useState } from "react";

// Dynamically import to avoid SSR issues
const ForceGraph2D = dynamic(() => import("react-force-graph-2d"), { ssr: false });

interface LiteraryGraphProps {
  selectedWriter: Writer | null;
  selectedWork: Work | null;
}

interface GraphNode {
  id: string;
  name: string;
  type: "writer" | "work";
  writerId?: number;
  workId?: number;
}

interface GraphLink {
  source: string;
  target: string;
  sentiment: boolean;
  quote: string;
  sourceRef: string;
}

export const LiteraryGraph: React.FC<LiteraryGraphProps> = ({
  selectedWriter,
  selectedWork,
}): React.JSX.Element => {
  const [nodes, setNodes] = useState<GraphNode[]>([]);
  const [links, setLinks] = useState<GraphLink[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  // ForceGraph2D ref type is not properly exported, using any is necessary
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const graphRef = useRef<any>(null);

  const loadGraphData = useCallback(async (): Promise<void> => {
    setIsLoading(true);
    setError(null);

    try {
      const [allWriters, allWorks, allOpinions] = await Promise.all([
        WriterService.list(1000, 0),
        WorkService.list(1000, 0),
        OpinionService.list(1000, 0),
      ]);

      const nodeMap = new Map<string, GraphNode>();
      const linkList: GraphLink[] = [];

      // Add all writers as nodes
      allWriters.forEach((writer) => {
        const nodeId = `writer-${writer.id}`;
        nodeMap.set(nodeId, {
          id: nodeId,
          name: writer.name,
          type: "writer",
          writerId: writer.id,
        });
      });

      // Add all works as nodes
      allWorks.forEach((work) => {
        const nodeId = `work-${work.id}`;
        nodeMap.set(nodeId, {
          id: nodeId,
          name: work.title,
          type: "work",
          workId: work.id,
        });
      });

      // Add links from opinions
      allOpinions.forEach((opinion) => {
        const writerNodeId = `writer-${opinion.writer_id}`;
        const workNodeId = `work-${opinion.work_id}`;

        if (nodeMap.has(writerNodeId) && nodeMap.has(workNodeId)) {
          linkList.push({
            source: writerNodeId,
            target: workNodeId,
            sentiment: opinion.sentiment,
            quote: opinion.quote,
            sourceRef: opinion.source,
          });
        }
      });

      // Filter based on selection
      let filteredNodes = Array.from(nodeMap.values());
      let filteredLinks = linkList;

      if (selectedWriter) {
        const writerNodeId = `writer-${selectedWriter.id}`;
        filteredLinks = filteredLinks.filter(
          (link) => link.source === writerNodeId || link.target === writerNodeId
        );
        // Add connected nodes
        const connectedNodeIds = new Set<string>([writerNodeId]);
        filteredLinks.forEach((link) => {
          connectedNodeIds.add(link.source);
          connectedNodeIds.add(link.target);
        });
        filteredNodes = filteredNodes.filter((node) => connectedNodeIds.has(node.id));
      } else if (selectedWork) {
        const workNodeId = `work-${selectedWork.id}`;
        filteredLinks = filteredLinks.filter(
          (link) => link.source === workNodeId || link.target === workNodeId
        );
        // Add connected nodes
        const connectedNodeIds = new Set<string>([workNodeId]);
        filteredLinks.forEach((link) => {
          connectedNodeIds.add(link.source);
          connectedNodeIds.add(link.target);
        });
        filteredNodes = filteredNodes.filter((node) => connectedNodeIds.has(node.id));
      }

      setNodes(filteredNodes);
      setLinks(filteredLinks);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load graph data");
      console.error("Error loading graph data:", err);
    } finally {
      setIsLoading(false);
    }
  }, [selectedWriter, selectedWork]);

  useEffect(() => {
    void loadGraphData();
  }, [loadGraphData]);

  if (isLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <div className="mx-auto h-8 w-8 animate-spin rounded-full border-4 border-gray-200 border-t-blue-600"></div>
          <p className="mt-4 text-gray-600">Loading graph...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <p className="text-red-600">Error: {error}</p>
          <button
            type="button"
            onClick={() => void loadGraphData()}
            className="mt-4 rounded-md bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  if (nodes.length === 0) {
    return (
      <div className="flex h-full items-center justify-center">
        <p className="text-gray-600">No data to display. Select a writer or work to explore.</p>
      </div>
    );
  }

  // Type-safe wrapper functions to avoid eslint-disable comments
  // The library uses 'any' types, but we know our data structure
  const getNodeLabel = (node: unknown): string => {
    const graphNode = node as GraphNode;
    return `${graphNode.name} (${graphNode.type})`;
  };

  const getNodeColor = (node: unknown): string => {
    const graphNode = node as GraphNode;
    return graphNode.type === "writer" ? "#3b82f6" : "#10b981";
  };

  const getLinkColor = (link: unknown): string => {
    const graphLink = link as GraphLink;
    return graphLink.sentiment ? "#10b981" : "#ef4444";
  };

  const getNodeVal = (node: unknown): number => {
    const graphNode = node as GraphNode;
    return graphNode.type === "writer" ? 8 : 6;
  };

  const handleLinkClick = (link: unknown): void => {
    const linkData = link as GraphLink;
    alert(`Opinion: ${linkData.quote}\n\nSource: ${linkData.sourceRef}`);
  };

  return (
    <ForceGraph2D
      ref={graphRef}
      graphData={{ nodes, links }}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      nodeLabel={getNodeLabel as any}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      nodeColor={getNodeColor as any}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      linkColor={getLinkColor as any}
      linkWidth={2}
      linkDirectionalArrowLength={6}
      linkDirectionalArrowRelPos={1}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      nodeVal={getNodeVal as any}
      onNodeClick={() => {
        // Could add node click handler here
      }}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      onLinkClick={handleLinkClick as any}
      cooldownTicks={100}
      onEngineStop={() => {
        if (graphRef.current) {
          graphRef.current.zoomToFit(400);
        }
      }}
    />
  );
};
