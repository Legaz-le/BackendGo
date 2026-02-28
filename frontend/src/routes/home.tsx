import { Link } from "@tanstack/react-router";

const HomePage = () => {
  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="text-center px-6">
        <h1 className="text-5xl font-bold text-gray-900 mb-4">Find Your Next Job</h1>
        <p className="text-xl text-gray-500 mb-8">Browse hundreds of job listings from top employers.</p>
        <Link
          to="/jobs"
          className="bg-blue-600 hover:bg-blue-700 text-white px-8 py-3 rounded-lg font-medium text-lg transition-colors"
        >
          Browse Jobs
        </Link>
      </div>
    </div>
  );
};

export default HomePage;
