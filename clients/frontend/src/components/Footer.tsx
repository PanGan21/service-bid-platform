export const Footer: React.FC = () => {
  return (
    <div
      className="fixed-bottom mt-3"
      style={{
        backgroundColor: "rgb(51,51,51)",
        paddingBottom: "5px",
        display: "flex",
        justifyContent: "center",
      }}
    >
      <div className="text-white">
        &copy; {new Date().getFullYear()} Copyright:{" "}
        <a
          className="text-white"
          href="https://github.com/PanGan21/"
          target="_blank"
          rel="noopener noreferrer"
        >
          Panagiotis Ganelis
        </a>
      </div>
    </div>
  );
};
