// src/components/Sidebar.jsx
"use client"
import Link from "next/link"; // Import Link for navigation
import { usePathname } from "next/navigation"; // Hook to check active page
import { LayoutDashboard, Wallet, Users, FileText, Settings } from "lucide-react";

const Sidebar = () => {
  const pathname = usePathname(); // Get current URL

  const menuItems = [
    { icon: LayoutDashboard, label: "Dashboard", href: "/" },
    { icon: Wallet, label: "Deposits & Funds", href: "/deposits" }, // New Link!
    { icon: Users, label: "Borrowers", href: "/borrowers" },
    { icon: FileText, label: "Reports", href: "/reports" },
    { icon: Settings, label: "Settings", href: "/settings" },
    {icon: Settings, label: "Loan Application", href: "/cooperative-info" },
  ];

  return (
    <aside className="w-64 bg-white dark:bg-[#111C44] hidden md:flex flex-col border-r border-gray-100 dark:border-none p-4 m-4 rounded-xl shadow-sm h-[calc(100vh-2rem)]">
      <div className="flex items-center gap-2 mb-8 px-2">
        <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white font-bold">C</div>
        <span className="text-xl font-bold text-[#1B254B] dark:text-white">Coop Manager</span>
      </div>

      <nav className="flex-1 space-y-2">
        {menuItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link 
              href={item.href} 
              key={item.label}
              className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-bold transition-colors ${
                isActive
                  ? "bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 border-r-4 border-blue-600"
                  : "text-gray-400 hover:bg-gray-50 dark:hover:bg-white/5 hover:text-[#1B254B] dark:hover:text-white"
              }`}
            >
              <item.icon className="w-5 h-5" />
              {item.label}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
};

export default Sidebar;