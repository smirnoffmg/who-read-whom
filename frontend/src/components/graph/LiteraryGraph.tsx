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

// Define canvas drawing functions outside component to avoid re-render issues
const nodeCanvasObject = (node: unknown, ctx: CanvasRenderingContext2D, globalScale: number): void => {
  const graphNode = node as GraphNode & {
    x: number;
    y: number;
  };
  const nodeSize = graphNode.type === "writer" ? 8 : 6;
  const nodeColor = graphNode.type === "writer" ? "#3b82f6" : "#10b981";

  // Save canvas state
  ctx.save();

  // Draw node circle
  ctx.beginPath();
  ctx.arc(graphNode.x, graphNode.y, nodeSize, 0, 2 * Math.PI);
  ctx.fillStyle = nodeColor;
  ctx.fill();
  ctx.closePath();

  // Draw text label below the node
  const label = graphNode.name;
  const fontSize = Math.max(10, 12 / globalScale);
  ctx.font = `${fontSize}px Sans-Serif`;
  const textWidth = ctx.measureText(label).width;
  const textHeight = fontSize;
  const padding = fontSize * 0.3;
  const bckgWidth = textWidth + padding * 2;
  const bckgHeight = textHeight + padding * 2;
  
  // Get canvas dimensions to prevent clipping
  const canvasWidth = ctx.canvas.width;
  const canvasHeight = ctx.canvas.height;
  
  // Calculate label position, adjusting to stay within canvas bounds
  let labelX = graphNode.x;
  let labelY = graphNode.y + nodeSize + 2;
  
  // Adjust X position if label would be clipped on left or right
  const halfWidth = bckgWidth / 2;
  if (labelX - halfWidth < 0) {
    labelX = halfWidth;
  } else if (labelX + halfWidth > canvasWidth) {
    labelX = canvasWidth - halfWidth;
  }
  
  // Adjust Y position if label would be clipped at bottom
  if (labelY + bckgHeight > canvasHeight) {
    labelY = graphNode.y - nodeSize - 2 - bckgHeight; // Place above node instead
  }

  // Draw background rectangle for text
  ctx.fillStyle = "rgba(255, 255, 255, 0.95)";
  ctx.fillRect(
    labelX - bckgWidth / 2,
    labelY,
    bckgWidth,
    bckgHeight
  );

  // Draw text
  ctx.textAlign = "center";
  ctx.textBaseline = "top";
  ctx.fillStyle = "#1f2937";
  ctx.fillText(label, labelX, labelY + padding);

  // Restore canvas state
  ctx.restore();
};

const nodePointerAreaPaint = (node: unknown, color: string, ctx: CanvasRenderingContext2D): void => {
  const graphNode = node as GraphNode & {
    x: number;
    y: number;
  };
  
  // Save canvas state
  ctx.save();
  
  const nodeSize = graphNode.type === "writer" ? 8 : 6;
  const label = graphNode.name;
  const fontSize = 12; // Use base font size for calculation
  ctx.font = `${fontSize}px Sans-Serif`;
  const textWidth = ctx.measureText(label).width;
  const padding = fontSize * 0.3;
  const textHeight = fontSize;
  const bckgWidth = textWidth + padding * 2;
  const bckgHeight = textHeight + padding * 2;
  const totalWidth = Math.max(nodeSize * 2, bckgWidth);
  const totalHeight = nodeSize + 2 + bckgHeight;

  // Draw hover area covering circle and text
  ctx.fillStyle = color;
  ctx.fillRect(
    graphNode.x - totalWidth / 2,
    graphNode.y - nodeSize,
    totalWidth,
    totalHeight
  );
  
  // Restore canvas state
  ctx.restore();
};

export const LiteraryGraph: React.FC<LiteraryGraphProps> = ({
  selectedWriter,
  selectedWork,
}): React.JSX.Element => {
  const [nodes, setNodes] = useState<GraphNode[]>([]);
  const [links, setLinks] = useState<GraphLink[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  // ForceGraph2D ref type is not properly exported, using any is necessary
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const graphRef = useRef<any>(null);

  const loadGraphData = useCallback(async (): Promise<void> => {
    setIsLoading(true);
    setError(null);

    try {
      // If nothing is selected, show empty graph
      if (!selectedWriter && !selectedWork) {
        setNodes([]);
        setLinks([]);
        setIsLoading(false);
        return;
      }

      const nodeMap = new Map<string, GraphNode>();
      const linkList: GraphLink[] = [];

      if (selectedWriter) {
        // Always add the selected writer node
        const writerNodeId = `writer-${selectedWriter.id}`;
        nodeMap.set(writerNodeId, {
          id: writerNodeId,
          name: selectedWriter.name,
          type: "writer",
          writerId: selectedWriter.id,
        });

        try {
          // Get opinions by this writer
          const opinions = await OpinionService.getByWriter(selectedWriter.id);
          
          if (opinions.length > 0) {
            // Get unique work IDs from opinions
            const workIds = [...new Set(opinions.map((o) => o.work_id))];
            
            // Load only the works that this writer has opinions about
            const workResults = await Promise.allSettled(
              workIds.map((workId) => WorkService.getById(workId))
            );

            // Add work nodes (filter out failed requests)
            workResults.forEach((result) => {
              if (result.status === "fulfilled" && result.value) {
                const work = result.value;
                const workNodeId = `work-${work.id}`;
                nodeMap.set(workNodeId, {
                  id: workNodeId,
                  name: work.title,
                  type: "work",
                  workId: work.id,
                });
              }
            });

            // Add links from opinions
            opinions.forEach((opinion) => {
              const workNodeId = `work-${opinion.work_id}`;
              if (nodeMap.has(workNodeId)) {
                linkList.push({
                  source: writerNodeId,
                  target: workNodeId,
                  sentiment: opinion.sentiment,
                  quote: opinion.quote,
                  sourceRef: opinion.source,
                });
              }
            });
          }
        } catch (err) {
          console.error("Error loading opinions for writer:", err);
          // Still show the writer node even if opinions fail to load
        }
      } else if (selectedWork) {
        // Validate work object
        if (!selectedWork.id || !selectedWork.title) {
          // eslint-disable-next-line no-console
          console.error("Invalid work object:", selectedWork);
          setError("Invalid work data: missing id or title");
          setIsLoading(false);
          return;
        }

        // Always add the selected work node
        const workNodeId = `work-${selectedWork.id}`;
        nodeMap.set(workNodeId, {
          id: workNodeId,
          name: selectedWork.title,
          type: "work",
          workId: selectedWork.id,
        });

        // eslint-disable-next-line no-console
        console.log("Selected work:", selectedWork);

        try {
          // Get opinions about this work
          const opinions = await OpinionService.getByWork(selectedWork.id);
          
          // eslint-disable-next-line no-console
          console.log(`Found ${opinions.length} opinions for work ${selectedWork.id}`);
          
          if (opinions.length > 0) {
            // Get unique writer IDs from opinions
            const writerIds = [...new Set(opinions.map((o) => o.writer_id))];
            
            // eslint-disable-next-line no-console
            console.log(`Loading ${writerIds.length} writers for work ${selectedWork.id}`);
            
            // Load only the writers who have opinions about this work
            const writerResults = await Promise.allSettled(
              writerIds.map((writerId) => WriterService.getById(writerId))
            );

            // Add writer nodes (filter out failed requests)
            let writersLoaded = 0;
            writerResults.forEach((result) => {
              if (result.status === "fulfilled" && result.value) {
                const writer = result.value;
                const writerNodeId = `writer-${writer.id}`;
                nodeMap.set(writerNodeId, {
                  id: writerNodeId,
                  name: writer.name,
                  type: "writer",
                  writerId: writer.id,
                });
                writersLoaded++;
              } else if (result.status === "rejected") {
                // eslint-disable-next-line no-console
                console.error("Failed to load writer:", result.reason);
              }
            });

            // eslint-disable-next-line no-console
            console.log(`Loaded ${writersLoaded} writers successfully`);

            // Add links from opinions
            opinions.forEach((opinion) => {
              const writerNodeId = `writer-${opinion.writer_id}`;
              if (nodeMap.has(writerNodeId)) {
                linkList.push({
                  source: writerNodeId,
                  target: workNodeId,
                  sentiment: opinion.sentiment,
                  quote: opinion.quote,
                  sourceRef: opinion.source,
                });
              }
            });
          } else {
            // eslint-disable-next-line no-console
            console.log(`No opinions found for work ${selectedWork.id}, showing work node only`);
          }
        } catch (err) {
          // eslint-disable-next-line no-console
          console.error("Error loading opinions for work:", err);
          // Still show the work node even if opinions fail to load
        }
      }

      const finalNodes = Array.from(nodeMap.values());
      setNodes(finalNodes);
      setLinks(linkList);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load graph data");
      console.error("Error loading graph data:", err);
    } finally {
      setIsLoading(false);
    }
  }, [selectedWriter, selectedWork]);

  useEffect(() => {
    void loadGraphData();
  }, [loadGraphData, selectedWriter, selectedWork]);

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

  // Show empty state only if nothing is selected
  // If something is selected but has no connections, still show the node
  if (nodes.length === 0) {
    if (!selectedWriter && !selectedWork) {
      return (
        <div className="flex h-full items-center justify-center">
          <p className="text-gray-600">
            Select a writer or work from the search bar to view the graph.
          </p>
        </div>
      );
    }
    // If something is selected but no nodes were loaded, show error message
    return (
      <div className="flex h-full items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600">
            {selectedWork
              ? `No data found for "${selectedWork.title}". The work may not have any opinions yet.`
              : `No data found for "${selectedWriter?.name}". The writer may not have any opinions yet.`}
          </p>
          {error && (
            <p className="mt-2 text-sm text-red-600">Error: {error}</p>
          )}
        </div>
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

  // Functions are defined outside component, so they're stable references

  // Create a key to force re-render when selection changes
  const graphKey = selectedWriter
    ? `writer-${selectedWriter.id}`
    : selectedWork
      ? `work-${selectedWork.id}`
      : "empty";

  return (
    <div
      onWheel={(e) => {
        e.preventDefault();
        e.stopPropagation();
      }}
      style={{ width: "100%", height: "100%" }}
    >
      <ForceGraph2D
        key={graphKey}
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
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      nodeCanvasObject={nodeCanvasObject as any}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      nodePointerAreaPaint={nodePointerAreaPaint as any}
      enableZoomInteraction={false}
      enablePanInteraction={false}
      minZoom={1}
      maxZoom={1}
      onZoom={() => {
        // Prevent zoom by resetting if it changes
        if (graphRef.current) {
          graphRef.current.zoom(1, 0);
        }
      }}
      onNodeClick={() => {
        // Could add node click handler here
      }}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      onLinkClick={handleLinkClick as any}
      cooldownTicks={100}
      onEngineStop={() => {
        if (graphRef.current && nodes.length > 0) {
          // Center and zoom to fit all nodes with extra padding for labels
          // Padding accounts for node size + label height (~50px)
          graphRef.current.zoomToFit(400, 50);
        }
      }}
      />
    </div>
  );
};
