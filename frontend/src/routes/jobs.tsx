import { useQuery } from "@tanstack/react-query";

type JobType = {
  id: number;
  title: string;
  description: string;
  location: string;
  min_salary: number;
  max_salary: number;
};

const JobsPage = () => {
  const { data, isLoading, isError } = useQuery<JobType[]>({
    queryKey: ["jobs"],
    queryFn: () => fetch("/jobs").then((res) => res.json()),
  });
  return (
    <div>
      {isLoading && <div>Loading...</div>}
      {isError && <div>Error loading</div>}
      {data && (
        <div>
          {data.map((item: JobType) => (
            <div key={item.id}>
              <h1>{item.title}</h1>
              <p>{item.description}</p>
              <p>{item.location}</p>
              <span>{item.max_salary}</span>
              <span>{item.min_salary}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default JobsPage;
