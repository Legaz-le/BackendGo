import { useNavigate } from "@tanstack/react-router";
import { Link } from "@tanstack/react-router";


const Navbar = () => {
  const navigate = useNavigate();
  const handleLogout = async () => {
    await fetch("/auth/logout", { method: "POST" });
    navigate({ to: "/login" });
  };
  return (
    <nav>
      <Link to="/">Home</Link>
      <Link to="/jobs">Jobs</Link>
      <Link to="/login">Login</Link>
      <Link to="/post-job">Post Job</Link>
      <button onClick={handleLogout}>Logout</button>
    </nav>
  );
};

export default Navbar;
