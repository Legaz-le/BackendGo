import {
  createRootRoute,
  createRoute,
  Outlet,
  Link,
} from "@tanstack/react-router";
import JobsPage from "./routes/jobs";
import LoginPage from "./routes/login";
import PostPage from "./routes/post-job";

const rootRoute = createRootRoute({
  component: () => (
    <>
      <nav>
        <Link to="/">Home</Link>
        <Link to="/jobs">Jobs</Link>
        <Link to="/login">Login</Link>
        <Link to="/post-job">Post Job</Link>
      </nav>
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
})

export const routeTree = rootRoute.addChildren([
  indexRoute,
  indexJob,
  indexLogin,
  indexPost
]);
