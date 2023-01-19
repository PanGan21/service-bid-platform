export default function authHeader() {
  const userStr = localStorage.getItem("user");
  let user = null;
  if (userStr) user = JSON.parse(userStr);

  if (user && user.accessToken) {
    return { Cookie: "s.id " + user.accessToken };
  } else {
    return { Cookie: "" };
  }
}
