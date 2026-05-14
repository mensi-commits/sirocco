import { Menu, Bell, Search } from "lucide-react";
import ThemeToggle from "../ThemeToggle";

type TopBarProps = {
    onToggleSidebar?: () => void;
};

export default function TopBar({ onToggleSidebar }: TopBarProps) {
    return (
        <header
            className="
        w-full h-16 
        flex items-center justify-between px-4 
        border-b transition-colors duration-300
        bg-white text-black border-gray-200
        dark:bg-gray-900 dark:text-white dark:border-gray-800
      "
        >
            {/* Left side */}
            <div className="flex items-center gap-3">
                {/* <button
                    onClick={onToggleSidebar}
                    className="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-800 transition"
                >
                    <Menu size={20} />
                </button> */}
            </div>

            {/* Center search */}
            <div
                className="
          hidden md:flex items-center px-3 py-2 rounded-lg w-1/3
          bg-gray-100 dark:bg-gray-800
        "
            >
                <Search size={18} className="text-gray-500 dark:text-gray-400" />
                <input
                    type="text"
                    placeholder="Search..."
                    className="
            bg-transparent outline-none text-sm ml-2 w-full
            text-black dark:text-white
          "
                />
            </div>

            {/* Right side */}
            <div className="flex items-center gap-4">
                <ThemeToggle />

                <button className="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-800 transition relative">
                    <Bell size={20} />
                    <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
                </button>

                <div className="w-8 h-8 rounded-full bg-blue-500 flex items-center justify-center text-sm font-bold text-white">
                    A
                </div>
            </div>
        </header>
    );
}