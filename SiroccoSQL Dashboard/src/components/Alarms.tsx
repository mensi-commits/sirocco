import { useState } from "react";
import { Bell, AlertTriangle, ShieldAlert, Info, CheckCircle } from "lucide-react";

type Severity = "info" | "warning" | "critical";

type Alarm = {
    id: string;
    title: string;
    message: string;
    severity: Severity;
    time: string;
    status: "active" | "acknowledged" | "resolved";
};

const initialAlarms: Alarm[] = [
    {
        id: "a1",
        title: "High Query Latency",
        message: "P95 latency exceeded 120ms on shard-2",
        severity: "critical",
        time: "2 min ago",
        status: "active",
    },
    {
        id: "a2",
        title: "Node Degraded",
        message: "node-02 experiencing high CPU load",
        severity: "warning",
        time: "10 min ago",
        status: "acknowledged",
    },
    {
        id: "a3",
        title: "Cache Miss Spike",
        message: "Cache hit rate dropped below 80%",
        severity: "info",
        time: "30 min ago",
        status: "resolved",
    },
];

export default function Alarms() {
    const [alarms, setAlarms] = useState(initialAlarms);

    const resolveAlarm = (id: string) => {
        setAlarms((prev) =>
            prev.map((a) =>
                a.id === id ? { ...a, status: "resolved" } : a
            )
        );
    };

    const ackAlarm = (id: string) => {
        setAlarms((prev) =>
            prev.map((a) =>
                a.id === id ? { ...a, status: "acknowledged" } : a
            )
        );
    };

    return (
        <div className="space-y-6">

            {/* HEADER */}
            <div>
                <h1 className="text-2xl font-bold">Alarms & Incidents</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Real-time system alerts & anomaly detection center
                </p>
            </div>

            {/* SUMMARY */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">

                <SummaryCard label="Critical" value="1" color="red" icon={<AlertTriangle />} />
                <SummaryCard label="Warnings" value="1" color="yellow" icon={<ShieldAlert />} />
                <SummaryCard label="Info" value="1" color="blue" icon={<Info />} />

            </div>

            {/* LIST */}
            <div className="space-y-3">

                {alarms.map((alarm) => (
                    <AlarmCard
                        key={alarm.id}
                        alarm={alarm}
                        onAck={ackAlarm}
                        onResolve={resolveAlarm}
                    />
                ))}

            </div>
        </div>
    );
}

/* ---------------- ALARM CARD ---------------- */

function AlarmCard({
    alarm,
    onAck,
    onResolve,
}: {
    alarm: Alarm;
    onAck: (id: string) => void;
    onResolve: (id: string) => void;
}) {
    const severityColor =
        alarm.severity === "critical"
            ? "text-red-500"
            : alarm.severity === "warning"
                ? "text-yellow-500"
                : "text-blue-500";

    const statusColor =
        alarm.status === "active"
            ? "text-red-500"
            : alarm.status === "acknowledged"
                ? "text-yellow-500"
                : "text-green-500";

    return (
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4">

            {/* TOP */}
            <div className="flex items-start justify-between">

                <div className="flex gap-3">

                    <Bell className={severityColor} />

                    <div>
                        <h3 className="font-semibold">{alarm.title}</h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                            {alarm.message}
                        </p>

                        <div className="text-xs mt-2 flex gap-3 text-gray-500">
                            <span>{alarm.time}</span>
                            <span className={severityColor}>{alarm.severity}</span>
                            <span className={statusColor}>{alarm.status}</span>
                        </div>
                    </div>
                </div>

                {/* ACTIONS */}
                <div className="flex gap-2">

                    {alarm.status !== "acknowledged" && (
                        <button
                            onClick={() => onAck(alarm.id)}
                            className="px-3 py-1 text-xs rounded bg-yellow-500 text-white hover:bg-yellow-600"
                        >
                            Acknowledge
                        </button>
                    )}

                    {alarm.status !== "resolved" && (
                        <button
                            onClick={() => onResolve(alarm.id)}
                            className="px-3 py-1 text-xs rounded bg-green-500 text-white hover:bg-green-600"
                        >
                            Resolve
                        </button>
                    )}

                </div>
            </div>
        </div>
    );
}

/* ---------------- SUMMARY CARD ---------------- */

function SummaryCard({
    label,
    value,
    color,
    icon,
}: {
    label: string;
    value: string;
    color: string;
    icon: React.ReactNode;
}) {
    const colorMap: Record<string, string> = {
        red: "text-red-500",
        yellow: "text-yellow-500",
        blue: "text-blue-500",
    };

    return (
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4 flex items-center justify-between">

            <div>
                <p className="text-sm text-gray-500">{label}</p>
                <p className="text-2xl font-bold">{value}</p>
            </div>

            <div className={colorMap[color]}>
                {icon}
            </div>

        </div>
    );
}