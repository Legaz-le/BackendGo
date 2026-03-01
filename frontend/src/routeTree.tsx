import { createRootRoute, createRoute, Outlet } from "@tanstack/react-router";
import JobsPage from "./routes/jobs";
import LoginPage from "./routes/login";
import PostPage from "./routes/post-job";
import HomePage from "./routes/home";
import Navbar from "./components/Navbar";
import RegisterPage from "./routes/register";
import JobDetailPage from "./routes/job-detail";

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
  component: HomePage,
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

const indexRegister = createRoute({
  getParentRoute: () => rootRoute,
  path: "/register",
  component: RegisterPage
})

const indexJobId = createRoute({
  getParentRoute: () => rootRoute,
  path: "/jobs/$id",
  component: JobDetailPage
})

export const routeTree = rootRoute.addChildren([
  indexRoute,
  indexJob,
  indexLogin,
  indexPost,
  indexRegister,
  indexJobId,
]);
