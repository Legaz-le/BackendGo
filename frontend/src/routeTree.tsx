import { createRootRoute, createRoute, Outlet } from "@tanstack/react-router";
import JobsPage from "./routes/jobs";
import LoginPage from "./routes/login";
import PostPage from "./routes/post-job";
import Navbar from "./components/Navbar";

const rootRoute = createRootRoute({
  component: () => (
    <>
      <Navbar />
      <Outlet />
    </>
  ),
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: () => <div>Job Board</div>,
});

const indexJob = createRoute({
  getParentRoute: () => rootRoute,
  path: "/jobs",
  component: JobsPage,
});

const indexLogin = createRoute({
  getParentRoute: () => rootRoute,
  path: "/login",
  component: LoginPage,
});

const indexPost = createRoute({
  getParentRoute: () => rootRoute,
  path: "/post-job",
  component: PostPage,
});

export const routeTree = rootRoute.addChildren([
  indexRoute,
  indexJob,
  indexLogin,
  indexPost,
]);
