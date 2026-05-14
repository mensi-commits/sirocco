import { Server, Database, ArrowRight } from "lucide-react";

type Node = {
    id: string;
    x: number;
    y: number;
    status: "healthy" | "degraded" | "down";
};

const nodes: Node[] = [
    { id: "Proxy", x: 50, y: 30, status: "healthy" },
    { id: "Shard-1", x: 20, y: 120, status: "healthy" },
    { id: "Shard-2", x: 50, y: 150, status: "degraded" },
    { id: "Shard-3", x: 80, y: 120, status: "healthy" },
    { id: "Replica-A", x: 20, y: 220, status: "healthy" },
    { id: "Replica-B", x: 80, y: 220, status: "down" },
];

const connections = [
    ["Proxy", "Shard-1"],
    ["Proxy", "Shard-2"],
    ["Proxy", "Shard-3"],
    ["Shard-1", "Replica-A"],
    ["Shard-3", "Replica-B"],
];

export default function Topology() {
    return (
        <div className="space-y-6">

            {/* HEADER */}
            <div>
                <h1 className="text-2xl font-bold">Cluster Topology</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Visual routing map of SiroccoSQL proxy cluster
                </p>
            </div>

            {/* GRAPH CONTAINER */}
            <div className="relative h-[500px] bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl overflow-hidden">

                {/* CONNECTION LINES */}
                <svg className="absolute w-full h-full">
                    {connections.map(([a, b], i) => {
                        const nodeA = nodes.find(n => n.id === a)!;
                        const nodeB = nodes.find(n => n.id === b)!;

                        return (
                            <line
                                key={i}
                                x1={`${nodeA.x}%`}
                                y1={`${nodeA.y}%`}
                                x2={`${nodeB.x}%`}
                                y2={`${nodeB.y}%`}
                                stroke="#60a5fa"
                                strokeWidth="2"
                                opacity="0.6"
                            />
                        );
                    })}
                </svg>

                {/* NODES */}
                {nodes.map((node) => (
                    <Node key={node.id} node={node} />
                ))}
            </div>
        </div>
    );
}

/* ---------------- NODE ---------------- */

function Node({ node }: { node: Node }) {
    const color =
        node.status === "healthy"
            ? "bg-green-500"
            : node.status === "degraded"
                ? "bg-yellow-500"
                : "bg-red-500";

    const icon =
        node.id === "Proxy" ? (
            <Server size={16} />
        ) : (
            <Database size={16} />
        );

    return (
        <div
            className="absolute flex items-center gap-2 px-3 py-2 rounded-lg bg-gray-100 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-sm shadow"
            style={{
                left: `${node.x}%`,
                top: `${node.y}%`,
                transform: "translate(-50%, -50%)",
            }}
        >
            <span className={`w-2 h-2 rounded-full ${color}`} />
            {icon}
            <span>{node.id}</span>
        </div>
    );
}