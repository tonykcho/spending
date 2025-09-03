export default function TopNavBar() {
    return (
        <nav className="bg-gray-800 p-4">
            <div className="mx-auto flex justify-between items-center">
                <div className="text-white text-lg font-bold">Spending Tracker</div>
                <ul className="flex space-x-4">
                    <li><a href="/" className="text-white hover:text-gray-300">Home</a></li>
                    <li><a href="/category" className="text-white hover:text-gray-300">Category</a></li>
                    {/* <li><a href="/about" className="text-white hover:text-gray-300">About</a></li> */}
                </ul>
            </div>
        </nav>
    );
}