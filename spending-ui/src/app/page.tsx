import Image from "next/image";

export default function Home() {
  return (
    <div className="flex flex-col flex-1 items-center justify-center">
      <h1 className="text-2xl font-bold">Welcome to the Spending Tracker</h1>
      <p className="mt-4">Track your expenses and manage your budget effectively.</p>
    </div>
  );
}
