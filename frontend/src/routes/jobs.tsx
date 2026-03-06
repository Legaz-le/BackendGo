import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { useState } from "react";

type JobType = {
  id: number;
  title: string;
  description: string;
  location: string;
  min_salary: number;
  max_salary: number;
};

type PaginatedJobs = {
  data: JobType[];
  page: number;
  limit: number;
  total: number;
};

const JobsPage = () => {
  const [page, setPage] = useState(1);
  const { data, isLoading, isError } = useQuery<PaginatedJobs>({
    queryKey: ["jobs", page],
    queryFn: () =>
      fetch(`/jobs?page=${page}&limit=10`).then((res) => res.json()),
  });

  return (
    <div className="min-h-screen bg-gray-100 py-10">
      <div className="max-w-4xl mx-auto px-6">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Job Listings</h1>
        <p className="text-gray-500 mb-8">Browse the latest opportunities</p>

        {isLoading && (
          <div className="text-center text-gray-500 py-20">Loading jobs...</div>
        )}
        {isError && (
          <div className="text-center text-red-500 py-20">
            Failed to load jobs. Please try again.
          </div>
        )}
        {data && data.data.length === 0 && (
          <div className="text-center text-gray-500 py-20">
            No jobs available yet.
          </div>
        )}
        {data &&
          data.data.map((item: JobType) => (
            <Link key={item.id} to={`/jobs/${item.id}`}>
              <div className="bg-white shadow-sm hover:shadow-md transition-shadow rounded-xl p-6 mb-4 border border-gray-200">
                <div className="flex items-start justify-between">
                  <h2 className="text-xl font-semibold text-gray-900">
                    {item.title}
                  </h2>
                  <span className="text-green-600 font-medium text-sm bg-green-50 px-3 py-1 rounded-full">
                    ${item.min_salary.toLocaleString()} – $
                    {item.max_salary.toLocaleString()}
                  </span>
                </div>
                <p className="text-gray-600 mt-2 text-sm">{item.description}</p>
                <div className="flex items-center gap-2 mt-4 text-gray-400 text-sm">
                  <span>📍</span>
                  <span>{item.location}</span>
                </div>
              </div>
            </Link>
          ))}
      </div>
      <div>
        <button disabled={page == 1} onClick={() => setPage((p) => p - 1)}>
          Previous
        </button>
        <button
          disabled={!data || page * 10 >= data.total}
          onClick={() => setPage((p) => p + 1)}
        >
          Next
        </button>
      </div>
    </div>
  );
};

export default JobsPage;
