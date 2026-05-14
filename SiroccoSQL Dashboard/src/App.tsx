import { useState } from "react";
import SideNavbar from "./components/SideNavbar";
import TopBar from "./components/TopBar";
import Dashboard from "./components/Dashboard";
function App() {
  const [sidebarOpen, setSidebarOpen] = useState(true);

  return (
    <div className="flex min-h-screen bg-white text-black dark:bg-gray-900 dark:text-white">
      <SideNavbar open={sidebarOpen} setOpen={setSidebarOpen} />

      <div className="flex-1 flex flex-col">
        <TopBar onToggleSidebar={() => setSidebarOpen(!sidebarOpen)} />

        <div>
          <Dashboard />
        </div>
      </div>
    </div>
  );
}

export default App;