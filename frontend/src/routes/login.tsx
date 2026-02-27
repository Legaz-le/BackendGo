import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";

type LoginCredentials = {
  email: string;
  password: string;
};

const LoginPage = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const navigation = useNavigate();

  const mutation = useMutation({
    onSuccess: () => {
      navigation({ to: "/jobs" });
    },
    mutationFn: (credentials: LoginCredentials) =>
      fetch("/auth/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(credentials),
      })
  });
  
  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        mutation.mutate({ email, password });
      }}
    >
      <label>Email</label>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <label>Password</label>
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <button type="submit" disabled={mutation.isPending}>
        Login
      </button>
    </form>
  );
};

export default LoginPage;
