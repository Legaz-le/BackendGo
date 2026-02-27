import {
  createRootRoute,
  createRoute,
  Outlet,
  Link,
} from "@tanstack/react-router";
import JobsPage from "./routes/jobs";
import LoginPage from "./routes/login";

const rootRoute = createRootRoute({
  component: () => (
    <>
      <nav>
        <Link to="/">Home</Link>
        <Link to="/jobs">Jobs</Link>
        <Link to="/login">Login</Link>
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

export const routeTree = rootRoute.addChildren([
  indexRoute,
  indexJob,
  indexLogin,
]);
