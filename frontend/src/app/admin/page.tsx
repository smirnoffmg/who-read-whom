export default function AdminHome(): React.JSX.Element {
  return (
    <div className="flex flex-col h-full">
      <div className="p-6 bg-white border-b border-gray-200">
        <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-sm text-gray-600 mt-1">Admin overview and quick access</p>
      </div>

      <div className="flex-1 p-6 overflow-auto">
        <div className="max-w-4xl">
          <p className="text-gray-600 mb-6">
            Use the sidebar to navigate to Writers, Works, or Opinions management pages.
          </p>
          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <h2 className="font-semibold text-blue-900 mb-2">Quick Tips</h2>
            <ul className="text-sm text-blue-800 space-y-1 list-disc list-inside">
              <li>All data tables support inline editing</li>
              <li>Use CSV import/export for bulk operations</li>
              <li>Click Edit on any row to modify data</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}
