import { createRootRoute, createRoute, Outlet } from "@tanstack/react-router";
import JobsPage from "./routes/jobs";

const rootRoute = createRootRoute({
  component: () => <Outlet />,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: () => <div>Job Board</div>,
});

const indexJob = createRoute({
  getParentRoute: () => rootRoute,
  path: "/jobs",
  component: JobsPage
})

export const routeTree = rootRoute.addChildren([indexRoute, indexJob]);
