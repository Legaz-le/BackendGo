import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";

type PostType = {
  title: string;
  description: string;
  location: string;
  min_salary: number;
  max_salary: number;
};

const PostPage = () => {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [location, setLocation] = useState("");
  const [minSalary, setMinSalary] = useState(0);
  const [maxSalary, setMaxSalary] = useState(0);

  const navigation = useNavigate();

  const mutate = useMutation({
    onSuccess: () => {
      navigation({ to: "/jobs" });
    },
    mutationFn: (credentials: PostType) =>
      fetch("/jobs", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(credentials),
      }),
  });
  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        mutate.mutate({ title, location, description, min_salary: minSalary, max_salary: maxSalary  });
      }}
    >
      <label>Title</label>
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
      />
      <label>Description</label>
      <input
        type="text"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
      />
      <label>Location</label>
      <input
        type="text"
        value={location}
        onChange={(e) => setLocation(e.target.value)}
      />
      <label>Min_Salary</label>
      <input
        type="number"
        value={minSalary}
        onChange={(e) => setMinSalary(Number(e.target.value))}
      />
      <label>Max_Salary</label>
      <input
        type="number"
        value={maxSalary}
        onChange={(e) => setMaxSalary(Number(e.target.value))}
      />
      <button type="submit">Post</button>
    </form>
  );
};


export default PostPage