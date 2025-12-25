"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import React from "react";

interface NavItem {
  href: string;
  label: string;
  description: string;
}

const navItems: NavItem[] = [
  {
    href: "/admin",
    label: "Dashboard",
    description: "Admin overview",
  },
  {
    href: "/admin/writers",
    label: "Writers",
    description: "Manage writers",
  },
  {
    href: "/admin/works",
    label: "Works",
    description: "Manage works",
  },
  {
    href: "/admin/opinions",
    label: "Opinions",
    description: "Manage opinions",
  },
];

export const AdminSidebar: React.FC = (): React.JSX.Element => {
  const pathname = usePathname();

  return (
    <aside className="w-64 bg-white border-r border-gray-200 h-full overflow-y-auto">
      <div className="p-6">
        <Link href="/admin" className="block mb-8">
          <h1 className="text-2xl font-bold text-gray-900">Admin Panel</h1>
          <p className="text-sm text-gray-600 mt-1">Management Console</p>
        </Link>

        <nav className="space-y-2">
          {navItems.map((item) => {
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`block px-4 py-3 rounded-lg transition-colors ${
                  isActive
                    ? "bg-blue-50 text-blue-700 font-medium border-l-4 border-blue-600"
                    : "text-gray-700 hover:bg-gray-50"
                }`}
              >
                <div className="font-medium">{item.label}</div>
                <div className="text-xs text-gray-500 mt-0.5">{item.description}</div>
              </Link>
            );
          })}
        </nav>

        <div className="mt-8 pt-8 border-t border-gray-200">
          <Link
            href="/"
            className="block px-4 py-2 text-sm text-gray-600 hover:text-gray-900 transition-colors"
          >
            ← Back to Main Site
          </Link>
        </div>
      </div>
    </aside>
  );
};

