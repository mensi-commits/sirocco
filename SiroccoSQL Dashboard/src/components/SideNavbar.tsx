import { useState } from "react";
import {
    Home,
    User,
    Settings,
    Folder,
    LogOut,
    Menu,
    X,
} from "lucide-react";

export default function SideNavbar() {
    const [open, setOpen] = useState(true);

    const menuItems = [
        { name: "Home", icon: <Home size={20} />, path: "/" },
        { name: "Profile", icon: <User size={20} />, path: "/profile" },
        { name: "Projects", icon: <Folder size={20} />, path: "/projects" },
        { name: "Settings", icon: <Settings size={20} />, path: "/settings" },
    ];

    return (
        <div className="flex">
            {/* Sidebar */}
            <div
                className={`
          h-screen transition-all duration-300 flex flex-col
          ${open ? "w-64" : "w-20"}

          /* LIGHT MODE */
          bg-white text-black border-r border-gray-200

          /* DARK MODE */
          dark:bg-gray-900 dark:text-white dark:border-gray-800
        `}
            >
                {/* Top Section */}
                <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-800">
                    {open && <h1 className="text-xl font-bold">Dashboard</h1>}

                    <button
                        onClick={() => setOpen(!open)}
                        className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition"
                    >
                        {open ? <X size={20} /> : <Menu size={20} />}
                    </button>
                </div>

                {/* Menu */}
                <div className="flex-1 p-2">
                    {menuItems.map((item, index) => (
                        <a
                            key={index}
                            href={item.path}
                            className="
                flex items-center gap-3 p-3 rounded-lg transition

                hover:bg-gray-100 dark:hover:bg-gray-800
              "
                        >
                            {item.icon}
                            {open && (
                                <span className="text-sm font-medium">
                                    {item.name}
                                </span>
                            )}
                        </a>
                    ))}
                </div>

                {/* Bottom Logout */}
                <div className="p-4 border-t border-gray-200 dark:border-gray-800">
                    <button className="
            flex items-center gap-3 w-full p-3 rounded-lg transition

            hover:bg-red-100 dark:hover:bg-red-600/30
          ">
                        <LogOut size={20} />
                        {open && <span className="text-sm font-medium">Logout</span>}
                    </button>
                </div>
            </div>
        </div>
    );
}