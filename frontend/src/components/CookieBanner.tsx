interface CookieBannerProps {
  onChoice: (choice: boolean) => void;
}

function CookieBanner(props: CookieBannerProps) {
  return (
    <div className="bg-blue-100 bg-opacity-75 fixed bottom-0 w-full border-t p-4 flex items-center justify-between">
      <p className="text-lg">
        This website uses cookies to ensure you get the best experience on our website.
      </p>
      <div>
        <button
          className="p-2 border rounded mr-2 bg-blue-500 hover:brightness-90 text-white font-semibold"
          onClick={() => props.onChoice(true)}
        >
          Accept
        </button>
        <button
          className="p-2 border border-blue-500 text-blue-500 rounded mr-2 hover:bg-blue-100"
          onClick={() => props.onChoice(false)}
        >
          Decline
        </button>
      </div>
    </div>
  );
}

export default CookieBanner;
