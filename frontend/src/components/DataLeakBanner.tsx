import { MdClose, MdOutlineWarningAmber } from "react-icons/md";
import { Link } from "react-router-dom";

interface DataLeakBannerProps {
  onClose: () => void;
}

function DataLeakBanner(props: DataLeakBannerProps) {
  return (
    <div className="bg-yellow-300 bg-opacity-75 p-2 fixed bottom-0 w-full">
      <div className="flex justify-between ml-5 font-semibold">
        <div>
          <div className="flex gap-2 items-center justify-center">
            <p className="text-lg font-bold">
              <MdOutlineWarningAmber className="inline-block mb-1 text-2xl" />
            </p>
            <p>
              We had a data leak on the <strong>31st of October 2024.</strong>
            </p>
            <p>
              If you have an account, please{" "}
              <Link
                to={"/change-password"}
                className="underline text-blue-600"
              >
                change your password
              </Link>{" "}
              immediately.
            </p>
            <p className="text-lg font-bold">
              <MdOutlineWarningAmber className="inline-block mb-1 text-2xl" />
            </p>
          </div>
          <div>
            <p className="text-xs">
              *If you changed your password after the 31st of October 2024, you are safe, and can
              dismiss this message.
            </p>
          </div>
        </div>
        <button
          onClick={props.onClose}
          className="text-lg font-bold"
        >
          <MdClose />
        </button>
      </div>
    </div>
  );
}

export default DataLeakBanner;
