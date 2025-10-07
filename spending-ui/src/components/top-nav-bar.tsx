import Link from "next/link";

export default function TopNavBar()
{
    return (
        <nav className="bg-gray-800 p-4">
            <div className="mx-auto flex justify-between items-center">
                <div className="text-white text-lg font-bold">Spending Tracker</div>
                <ul className="flex space-x-4">
                    <li>
                        <Link href="/category" className="text-white hover:text-gray-300">Category</Link>
                    </li>
                    <li>
                        <Link href="/spending" className="text-white hover:text-gray-300">Spending</Link>
                    </li>
                </ul>
            </div>
        </nav>
    );
}