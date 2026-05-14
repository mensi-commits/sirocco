import { Activity, Cpu, Database, Zap, Clock } from "lucide-react";

export default function Performance() {
    return (
        <div className="space-y-6">

            {/* HEADER */}
            <div>
                <h1 className="text-2xl font-bold">Performance</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    System-wide throughput, latency & resource utilization
                </p>
            </div>

            {/* KPI GRID */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">

                <MetricCard
                    title="Queries / sec"
                    value="14,892"
                    icon={<Zap />}
                    trend="+6.3%"
                />

                <MetricCard
                    title="P95 Latency"
                    value="8.4 ms"
                    icon={<Clock />}
                    trend="-1.8%"
                />

                <MetricCard
                    title="Cache Hit Rate"
                    value="92.7%"
                    icon={<Database />}
                    trend="+2.1%"
                />

                <MetricCard
                    title="CPU Usage"
                    value="68%"
                    icon={<Cpu />}
                    trend="stable"
                />
            </div>

            {/* MAIN GRID */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">

                {/* LATENCY CHART */}
                <div className="lg:col-span-2 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4">
                    <h2 className="font-semibold flex items-center gap-2 mb-4">
                        <Activity size={18} /> Latency Trend
                    </h2>

                    <FakeChart />
                </div>

                {/* BREAKDOWN */}
                <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4">
                    <h2 className="font-semibold mb-4">Request Breakdown</h2>

                    <Breakdown label="Cached Queries" value={62} />
                    <Breakdown label="Routed Queries" value={28} />
                    <Breakdown label="Slow Queries" value={7} />
                    <Breakdown label="Failed Queries" value={3} />
                </div>
            </div>
        </div>
    );
}

/* ---------------- KPI CARD ---------------- */

function MetricCard({
    title,
    value,
    icon,
    trend,
}: {
    title: string;
    value: string;
    icon: React.ReactNode;
    trend: string;
}) {
    return (
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4">
            <div className="flex items-center justify-between text-gray-500 dark:text-gray-400 text-sm">
                {title}
                <div className="text-blue-500">{icon}</div>
            </div>

            <div className="text-2xl font-bold mt-2">{value}</div>

            <div className="text-xs text-gray-500 mt-1">{trend}</div>
        </div>
    );
}

/* ---------------- FAKE CHART ---------------- */

function FakeChart() {
    const bars = [20, 40, 35, 60, 80, 55, 90, 70, 85, 95];

    return (
        <div className="flex items-end gap-2 h-48">
            {bars.map((h, i) => (
                <div
                    key={i}
                    className="w-6 bg-blue-500 rounded-t"
                    style={{ height: `${h}%` }}
                />
            ))}
        </div>
    );
}

/* ---------------- BREAKDOWN ---------------- */

function Breakdown({
    label,
    value,
}: {
    label: string;
    value: number;
}) {
    return (
        <div className="mb-3">
            <div className="flex justify-between text-sm mb-1">
                <span>{label}</span>
                <span>{value}%</span>
            </div>

            <div className="w-full h-2 bg-gray-200 dark:bg-gray-800 rounded">
                <div
                    className="h-2 bg-blue-500 rounded"
                    style={{ width: `${value}%` }}
                />
            </div>
        </div>
    );
}