import { SpinnerInfinity } from "spinners-react";

function LoadingSpinner({
  size,
  thickness,
  speed
}: {
  size?: number;
  thickness?: number;
  speed?: number;
}) {
  return (
    <SpinnerInfinity
      size={size}
      thickness={thickness}
      speed={speed}
      color="rgba(59, 130, 246, 1)"
      secondaryColor="rgba(0, 0, 0, 0.44)"
    />
  );
}

export default LoadingSpinner;
