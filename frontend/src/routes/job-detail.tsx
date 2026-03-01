import { useQuery } from "@tanstack/react-query";
import { useParams, Link } from "@tanstack/react-router";

type JobType = {
  id: number;
  title: string;
  description: string;
  location: string;
  min_salary: number;
  max_salary: number;
};

const JobDetailPage = () => {
  const { id } = useParams({ strict: false });

  const { data, isLoading, isError } = useQuery<JobType>({
    queryKey: ["job", id],
    queryFn: () => fetch(`/jobs/${id}`).then((res) => res.json()),
  });

  if (isLoading) return <div className="min-h-screen flex items-center justify-center text-gray-500">Loading...</div>;
  if (isError || !data) return <div className="min-h-screen flex items-center justify-center text-red-500">Job not found.</div>;

  return (
    <div className="min-h-screen bg-gray-100 py-10">
      <div className="max-w-3xl mx-auto px-6">
        <Link to="/jobs" className="text-blue-600 hover:underline text-sm mb-6 inline-block">
          &larr; Back to Jobs
        </Link>
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8">
          <div className="flex items-start justify-between mb-4">
            <h1 className="text-3xl font-bold text-gray-900">{data.title}</h1>
            <span className="text-green-600 font-medium text-sm bg-green-50 px-3 py-1 rounded-full">
              ${data.min_salary.toLocaleString()} – ${data.max_salary.toLocaleString()}
            </span>
          </div>
          <div className="flex items-center gap-2 text-gray-400 text-sm mb-6">
            <span>📍</span>
            <span>{data.location}</span>
          </div>
          <p className="text-gray-700 leading-relaxed">{data.description}</p>
        </div>
      </div>
    </div>
  );
};

export default JobDetailPage;
