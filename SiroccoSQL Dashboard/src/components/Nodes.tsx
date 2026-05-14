import { Server, Cpu, Activity, AlertCircle } from "lucide-react";

type NodeStatus = "healthy" | "degraded" | "down";

type Node = {
    id: string;
    region: string;
    load: number;
    latency: number;
    shards: number;
    status: NodeStatus;
};

const nodes: Node[] = [
    { id: "node-01", region: "eu-central", load: 32, latency: 4.2, shards: 12, status: "healthy" },
    { id: "node-02", region: "eu-west", load: 78, latency: 12.5, shards: 18, status: "degraded" },
    { id: "node-03", region: "eu-south", load: 15, latency: 3.1, shards: 6, status: "healthy" },
    { id: "node-04", region: "eu-north", load: 0, latency: 0, shards: 0, status: "down" },
];

export default function Nodes() {
    return (
        <div className="space-y-6">

            {/* HEADER */}
            <div>
                <h1 className="text-2xl font-bold">Cluster Nodes</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Real-time distributed node health & shard allocation
                </p>
            </div>

            {/* GRID */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">

                {nodes.map((node) => (
                    <NodeCard key={node.id} node={node} />
                ))}
            </div>
        </div>
    );
}

/* ---------------- NODE CARD ---------------- */

function NodeCard({ node }: { node: Node }) {
    const statusColor =
        node.status === "healthy"
            ? "text-green-500"
            : node.status === "degraded"
                ? "text-yellow-500"
                : "text-red-500";

    return (
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4 space-y-4 hover:scale-[1.01] transition">

            {/* TOP */}
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                    <Server className="text-blue-500" size={18} />
                    <span className="font-semibold">{node.id}</span>
                </div>

                <span className={`text-sm font-medium ${statusColor}`}>
                    {node.status}
                </span>
            </div>

            {/* REGION */}
            <p className="text-sm text-gray-500 dark:text-gray-400">
                Region: {node.region}
            </p>

            {/* METRICS */}
            <div className="grid grid-cols-3 gap-2 text-sm">

                <Metric label="Load" value={`${node.load}%`} icon={<Cpu size={16} />} />
                <Metric label="Latency" value={`${node.latency}ms`} icon={<Activity size={16} />} />
                <Metric label="Shards" value={`${node.shards}`} icon={<Server size={16} />} />

            </div>

            {/* ACTIONS */}
            <div className="flex gap-2 pt-2">
                <button className="px-3 py-1 text-xs rounded bg-blue-500 text-white hover:bg-blue-600">
                    Inspect
                </button>

                <button className="px-3 py-1 text-xs rounded bg-yellow-500 text-white hover:bg-yellow-600">
                    Drain
                </button>

                <button className="px-3 py-1 text-xs rounded bg-red-500 text-white hover:bg-red-600">
                    Restart
                </button>
            </div>
        </div>
    );
}

/* ---------------- METRIC ---------------- */

function Metric({
    label,
    value,
    icon,
}: {
    label: string;
    value: string;
    icon: React.ReactNode;
}) {
    return (
        <div className="bg-gray-50 dark:bg-gray-800 p-2 rounded-lg">
            <div className="flex items-center gap-1 text-gray-500 dark:text-gray-400 text-xs">
                {icon} {label}
            </div>
            <div className="font-semibold">{value}</div>
        </div>
    );
}