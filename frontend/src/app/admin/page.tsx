import Link from "next/link";

export default function AdminHome(): React.JSX.Element {
  return (
    <div className="flex min-h-screen flex-col bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Admin Panel</h1>
              <p className="mt-2 text-sm text-gray-600">Manage writers, works, and opinions</p>
            </div>
            <Link
              href="/"
              className="rounded-md bg-gray-600 px-4 py-2 text-sm font-medium text-white hover:bg-gray-700"
            >
              View Main Page
            </Link>
          </div>
        </div>
      </header>

      <main className="mx-auto w-full max-w-7xl flex-1 px-4 py-8 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-3">
          <Link
            href="/admin/writers"
            className="group rounded-lg bg-white p-6 shadow-md transition-shadow hover:shadow-lg"
          >
            <h2 className="mb-2 text-xl font-semibold text-gray-900 group-hover:text-blue-600">
              Writers
            </h2>
            <p className="text-gray-600">Manage literary writers and their information</p>
          </Link>

          <Link
            href="/admin/works"
            className="group rounded-lg bg-white p-6 shadow-md transition-shadow hover:shadow-lg"
          >
            <h2 className="mb-2 text-xl font-semibold text-gray-900 group-hover:text-blue-600">
              Works
            </h2>
            <p className="text-gray-600">Track literary works and their authors</p>
          </Link>

          <Link
            href="/admin/opinions"
            className="group rounded-lg bg-white p-6 shadow-md transition-shadow hover:shadow-lg"
          >
            <h2 className="mb-2 text-xl font-semibold text-gray-900 group-hover:text-blue-600">
              Opinions
            </h2>
            <p className="text-gray-600">Document opinions with quotes and sources</p>
          </Link>
        </div>
      </main>
    </div>
  );
}
