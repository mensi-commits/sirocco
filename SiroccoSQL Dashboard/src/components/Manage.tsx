import { useState } from "react";
import { Users, Settings, Shield, Server, Key, RefreshCw } from "lucide-react";

type Tab = "users" | "cluster" | "routing" | "security";

export default function Manage() {
    const [tab, setTab] = useState<Tab>("users");

    return (
        <div className="space-y-6">

            {/* HEADER */}
            <div>
                <h1 className="text-2xl font-bold">Manage System</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Admin controls for SiroccoSQL proxy infrastructure
                </p>
            </div>

            {/* NAV TABS */}
            <div className="flex gap-2 flex-wrap">
                <TabButton active={tab === "users"} onClick={() => setTab("users")} icon={<Users size={16} />}>
                    Users
                </TabButton>

                <TabButton active={tab === "cluster"} onClick={() => setTab("cluster")} icon={<Server size={16} />}>
                    Cluster
                </TabButton>

                <TabButton active={tab === "routing"} onClick={() => setTab("routing")} icon={<RefreshCw size={16} />}>
                    Routing
                </TabButton>

                <TabButton active={tab === "security"} onClick={() => setTab("security")} icon={<Shield size={16} />}>
                    Security
                </TabButton>
            </div>

            {/* CONTENT */}
            <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-6">

                {tab === "users" && <UsersPanel />}
                {tab === "cluster" && <ClusterPanel />}
                {tab === "routing" && <RoutingPanel />}
                {tab === "security" && <SecurityPanel />}

            </div>
        </div>
    );
}

/* ---------------- TAB BUTTON ---------------- */

function TabButton({
    children,
    active,
    onClick,
    icon,
}: {
    children: React.ReactNode;
    active: boolean;
    onClick: () => void;
    icon: React.ReactNode;
}) {
    return (
        <button
            onClick={onClick}
            className={`
        flex items-center gap-2 px-4 py-2 rounded-lg text-sm transition

        ${active
                    ? "bg-blue-500 text-white"
                    : "bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700"
                }
      `}
        >
            {icon}
            {children}
        </button>
    );
}

/* ---------------- USERS ---------------- */

function UsersPanel() {
    return (
        <div className="space-y-4">
            <h2 className="font-semibold flex items-center gap-2">
                <Users size={18} /> User Management
            </h2>

            <UserRow name="admin@sirocco.io" role="Admin" />
            <UserRow name="dev@sirocco.io" role="Developer" />
            <UserRow name="viewer@sirocco.io" role="Read-only" />
        </div>
    );
}

function UserRow({ name, role }: { name: string; role: string }) {
    return (
        <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
            <div>
                <p className="font-medium">{name}</p>
                <p className="text-xs text-gray-500">{role}</p>
            </div>

            <button className="px-3 py-1 text-xs rounded bg-red-500 text-white hover:bg-red-600">
                Revoke
            </button>
        </div>
    );
}

/* ---------------- CLUSTER ---------------- */

function ClusterPanel() {
    return (
        <div className="space-y-3">
            <h2 className="font-semibold flex items-center gap-2">
                <Server size={18} /> Cluster Control
            </h2>

            <Action label="Scale Up Nodes" />
            <Action label="Drain Cluster" />
            <Action label="Restart Proxy Layer" />
        </div>
    );
}

/* ---------------- ROUTING ---------------- */

function RoutingPanel() {
    return (
        <div className="space-y-3">
            <h2 className="font-semibold flex items-center gap-2">
                <RefreshCw size={18} /> Routing Engine
            </h2>

            <Action label="Rebalance Shards" />
            <Action label="Clear Routing Cache" />
            <Action label="Optimize Query Paths" />
        </div>
    );
}

/* ---------------- SECURITY ---------------- */

function SecurityPanel() {
    return (
        <div className="space-y-3">
            <h2 className="font-semibold flex items-center gap-2">
                <Shield size={18} /> Security Controls
            </h2>

            <Action label="Rotate API Keys" />
            <Action label="Enable Query Firewall" />
            <Action label="Audit Logs Export" />
        </div>
    );
}

/* ---------------- ACTION ---------------- */

function Action({ label }: { label: string }) {
    return (
        <button className="w-full flex items-center justify-between p-3 rounded-lg bg-gray-50 dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 transition">
            <span>{label}</span>
            <Key size={16} className="text-blue-500" />
        </button>
    );
}