import { useNavigate } from "@tanstack/react-router";
import { Link } from "@tanstack/react-router";

const Navbar = () => {
  const navigate = useNavigate();
  const handleLogout = async () => {
    await fetch("/auth/logout", { method: "POST" });
    navigate({ to: "/login" });
  };
  return (
    <nav className="bg-gray-900 text-white px-6 py-4 flex items-center gap-6">
      <Link className="hover:text-blue-400 transition-colors" to="/">
        Home
      </Link>
      <Link className="hover:text-blue-400 transition-colors" to="/jobs">
        Jobs
      </Link>
      <Link className="hover:text-blue-400 transition-colors" to="/login">
        Login
      </Link>
      <Link className="hover:text-blue-400 transition-colors" to="/post-job">
        Post Job
      </Link>
      <button
        className="ml-auto bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded transition-colors"
        onClick={handleLogout}
      >
        Logout
      </button>
    </nav>
  );
};

export default Navbar;
