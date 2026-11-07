import React from 'react';
import { createRoot } from 'react-dom/client';
import './style.css';
function App() { return <main className="min-h-screen bg-slate-950 p-8 text-slate-100"><h1 className="text-3xl font-semibold">Weather Dashboard</h1><p className="mt-3 text-slate-400">Project setup is ready.</p></main>; }
createRoot(document.getElementById('root')!).render(<React.StrictMode><App /></React.StrictMode>);
