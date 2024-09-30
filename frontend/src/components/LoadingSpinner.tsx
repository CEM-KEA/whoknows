import { SpinnerInfinity } from "spinners-react";

interface LoadingSpinnerProps {
  size?: number;
  thickness?: number;
  speed?: number;
}

// This gives a "defaultProps will be removed..." warning,
// This is a known issue with spinners-react, and it is not a problem.
// The warning can be ignored.
function LoadingSpinner({ size, thickness, speed }: Readonly<LoadingSpinnerProps>) {
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
