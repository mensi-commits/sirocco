import { Database, Download, RotateCcw, ShieldCheck, Clock } from "lucide-react";

type BackupStatus = "completed" | "running" | "failed";

type BackupItem = {
    id: string;
    name: string;
    size: string;
    createdAt: string;
    status: BackupStatus;
};

const backups: BackupItem[] = [
    {
        id: "b1",
        name: "snapshot-prod-001",
        size: "2.4 GB",
        createdAt: "2026-05-14 10:12",
        status: "completed",
    },
    {
        id: "b2",
        name: "snapshot-prod-002",
        size: "2.5 GB",
        createdAt: "2026-05-13 10:12",
        status: "completed",
    },
    {
        id: "b3",
        name: "snapshot-prod-003",
        size: "2.6 GB",
        createdAt: "2026-05-12 10:12",
        status: "running",
    },
];

export default function Backup() {
    return (
        <div className="space-y-6">

            {/* HEADER */}
            <div>
                <h1 className="text-2xl font-bold">Backup & Recovery</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Cluster snapshots, restore points & disaster recovery
                </p>
            </div>

            {/* ACTION BAR */}
            <div className="flex flex-wrap gap-3">

                <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-blue-500 text-white hover:bg-blue-600">
                    <Database size={16} />
                    Create Snapshot
                </button>

                <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-gray-100 dark:bg-gray-800">
                    <ShieldCheck size={16} />
                    Enable Auto Backup
                </button>

                <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-gray-100 dark:bg-gray-800">
                    <RotateCcw size={16} />
                    Restore Cluster
                </button>
            </div>

            {/* GRID */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">

                {/* BACKUP LIST */}
                <div className="lg:col-span-2 space-y-3">

                    {backups.map((b) => (
                        <BackupCard key={b.id} backup={b} />
                    ))}

                </div>

                {/* INFO PANEL */}
                <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4 space-y-4">

                    <h2 className="font-semibold">Backup Status</h2>

                    <Stat label="Total Backups" value="18" />
                    <Stat label="Storage Used" value="42.3 GB" />
                    <Stat label="Retention Policy" value="7 days" />

                    <div className="pt-2 text-sm text-green-500 flex items-center gap-2">
                        <Clock size={16} />
                        Last backup: 12 min ago
                    </div>
                </div>

            </div>
        </div>
    );
}

/* ---------------- BACKUP CARD ---------------- */

function BackupCard({ backup }: { backup: BackupItem }) {
    const statusColor =
        backup.status === "completed"
            ? "text-green-500"
            : backup.status === "running"
                ? "text-yellow-500"
                : "text-red-500";

    return (
        <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-4 flex items-center justify-between hover:shadow transition">

            {/* LEFT */}
            <div className="flex items-center gap-3">
                <Database className="text-blue-500" size={18} />

                <div>
                    <p className="font-semibold">{backup.name}</p>

                    <div className="text-xs text-gray-500 dark:text-gray-400 flex gap-3">
                        <span>{backup.size}</span>
                        <span>{backup.createdAt}</span>
                        <span className={statusColor}>{backup.status}</span>
                    </div>
                </div>
            </div>

            {/* ACTIONS */}
            <div className="flex gap-2">
                <button className="px-3 py-1 text-xs rounded bg-blue-500 text-white hover:bg-blue-600">
                    Restore
                </button>

                <button className="px-3 py-1 text-xs rounded bg-gray-100 dark:bg-gray-800">
                    Download
                </button>
            </div>
        </div>
    );
}

/* ---------------- STATS ---------------- */

function Stat({ label, value }: { label: string; value: string }) {
    return (
        <div className="flex justify-between text-sm">
            <span className="text-gray-500 dark:text-gray-400">{label}</span>
            <span className="font-semibold">{value}</span>
        </div>
    );
}