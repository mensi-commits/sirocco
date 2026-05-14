import { useState } from "react";

type Tab =
    | "overview"
    | "nodes"
    | "topology"
    | "query"
    | "performance"
    | "backup"
    | "manage"
    | "security"
    | "alarms"
    | "logs";

const tabs: { key: Tab; label: string }[] = [
    { key: "overview", label: "Overview" },
    { key: "nodes", label: "Nodes" },
    { key: "topology", label: "Topology" },
    { key: "query", label: "Query Monitor" },
    { key: "performance", label: "Performance" },
    { key: "backup", label: "Backup" },
    { key: "manage", label: "Manage" },
    { key: "security", label: "Security" },
    { key: "alarms", label: "Alarms" },
    { key: "logs", label: "Logs" },
];

import Overview from "./Overview";
import Nodes from "./Nodes";
import Topology from "./Topology";
import QueryMonitor from "./QueryMonitor";
import Performance from "./Performance";
import Backup from "./Backup";

import Manage from "./Manage";
import Security from "./Security";

import Alarms from "./Alarms";
import Logs from "./Logs";


export default function Dashboard() {
    const [activeTab, setActiveTab] = useState<Tab>("overview");

    return (
        <div className="flex flex-col h-screen">
            {/* NAVBAR */}
            <div className="flex gap-2 p-2 border-b bg-white dark:bg-gray-900 dark:border-gray-800 overflow-x-auto">
                {tabs.map((tab) => (
                    <button
                        key={tab.key}
                        onClick={() => setActiveTab(tab.key)}
                        className={`
              px-4 py-2 rounded-lg text-sm whitespace-nowrap transition
              ${activeTab === tab.key
                                ? "bg-blue-500 text-white"
                                : "hover:bg-gray-100 dark:hover:bg-gray-800"
                            }
            `}
                    >
                        {tab.label}
                    </button>
                ))}
            </div>

            {/* CONTENT */}
            <div className="flex-1 p-6 bg-gray-50 dark:bg-gray-950 text-black dark:text-white">
                {activeTab === "overview" && <Overview />}
                {activeTab === "nodes" && <Nodes />}
                {activeTab === "topology" && <Topology />}
                {activeTab === "query" && <QueryMonitor />}
                {activeTab === "performance" && <Performance />}
                {activeTab === "backup" && <Backup />}
                {activeTab === "manage" && <Manage />}
                {activeTab === "security" && <Security />}
                {activeTab === "alarms" && <Alarms />}
                {activeTab === "logs" && <Logs />}
            </div>
        </div>
    );
}