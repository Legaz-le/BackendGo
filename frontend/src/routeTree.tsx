import { createRootRoute, createRoute, Outlet } from "@tanstack/react-router";

const rootRoute = createRootRoute({
  component: () => <Outlet />,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: () => <div>Job Board</div>,
});

export const routeTree = rootRoute.addChildren([indexRoute]);
