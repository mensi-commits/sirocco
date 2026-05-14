import { useState } from "react";
import { Shield, Key, Lock, Eye, AlertTriangle, RefreshCw } from "lucide-react";

type Tab = "overview" | "keys" | "firewall" | "audit";

export default function Security() {
    const [tab, setTab] = useState<Tab>("overview");

    return (
        <div className="space-y-6">

            {/* HEADER */}
            <div>
                <h1 className="text-2xl font-bold">Security Center</h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    Access control, encryption, firewall & audit monitoring
                </p>
            </div>

            {/* NAV */}
            <div className="flex gap-2 flex-wrap">

                <TabButton active={tab === "overview"} onClick={() => setTab("overview")} icon={<Shield size={16} />}>
                    Overview
                </TabButton>

                <TabButton active={tab === "keys"} onClick={() => setTab("keys")} icon={<Key size={16} />}>
                    API Keys
                </TabButton>

                <TabButton active={tab === "firewall"} onClick={() => setTab("firewall")} icon={<Lock size={16} />}>
                    Firewall
                </TabButton>

                <TabButton active={tab === "audit"} onClick={() => setTab("audit")} icon={<Eye size={16} />}>
                    Audit Logs
                </TabButton>

            </div>

            {/* CONTENT */}
            <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-xl p-6">

                {tab === "overview" && <Overview />}
                {tab === "keys" && <Keys />}
                {tab === "firewall" && <Firewall />}
                {tab === "audit" && <Audit />}

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

/* ---------------- OVERVIEW ---------------- */

function Overview() {
    return (
        <div className="space-y-4">

            <h2 className="font-semibold flex items-center gap-2">
                <Shield size={18} /> System Security Status
            </h2>

            <Status label="Encryption (TLS 1.3)" status="active" />
            <Status label="Query Firewall" status="active" />
            <Status label="Rate Limiting" status="active" />
            <Status label="Intrusion Detection" status="warning" />

        </div>
    );
}

/* ---------------- API KEYS ---------------- */

function Keys() {
    return (
        <div className="space-y-3">

            <h2 className="font-semibold flex items-center gap-2">
                <Key size={18} /> API Keys
            </h2>

            <KeyRow name="prod-key-1" status="active" />
            <KeyRow name="dev-key-2" status="active" />
            <KeyRow name="legacy-key-3" status="revoked" />

            <button className="mt-3 px-4 py-2 rounded-lg bg-blue-500 text-white hover:bg-blue-600">
                Generate New Key
            </button>
        </div>
    );
}

function KeyRow({ name, status }: { name: string; status: string }) {
    return (
        <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
            <span>{name}</span>
            <span className={`text-xs ${status === "active" ? "text-green-500" : "text-red-500"}`}>
                {status}
            </span>
        </div>
    );
}

/* ---------------- FIREWALL ---------------- */

function Firewall() {
    return (
        <div className="space-y-3">

            <h2 className="font-semibold flex items-center gap-2">
                <Lock size={18} /> Query Firewall Rules
            </h2>

            <Rule label="Block DROP statements in production" />
            <Rule label="Limit SELECT * on large tables" />
            <Rule label="Detect SQL injection patterns" />

        </div>
    );
}

function Rule({ label }: { label: string }) {
    return (
        <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
            <span>{label}</span>
            <span className="text-green-500 text-xs">enabled</span>
        </div>
    );
}

/* ---------------- AUDIT ---------------- */

function Audit() {
    return (
        <div className="space-y-3">

            <h2 className="font-semibold flex items-center gap-2">
                <Eye size={18} /> Audit Logs
            </h2>

            <AuditRow action="User admin executed DROP TABLE logs" type="critical" />
            <AuditRow action="Dev queried 1.2M rows" type="warning" />
            <AuditRow action="System rotated API keys" type="info" />

        </div>
    );
}

function AuditRow({ action, type }: { action: string; type: string }) {
    const color =
        type === "critical"
            ? "text-red-500"
            : type === "warning"
                ? "text-yellow-500"
                : "text-blue-500";

    return (
        <div className="p-3 bg-gray-50 dark:bg-gray-800 rounded-lg flex justify-between">
            <span>{action}</span>
            <span className={`text-xs ${color}`}>{type}</span>
        </div>
    );
}

/* ---------------- STATUS ---------------- */

function Status({ label, status }: { label: string; status: "active" | "warning" }) {
    return (
        <div className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
            <span>{label}</span>
            <span className={`text-xs ${status === "active" ? "text-green-500" : "text-yellow-500"}`}>
                {status}
            </span>
        </div>
    );
}