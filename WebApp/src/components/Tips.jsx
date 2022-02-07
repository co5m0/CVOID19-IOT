import React from "react";
import quotes from "../const/quotes";
import { useState, useEffect } from "react";

function Tips(props) {
  const [quote, setQuote] = useState(
    quotes[Math.floor(Math.random() * Math.floor(quotes.length))]
  );

  useEffect(() => {
    const i_id = setInterval(() => {
      var q = quotes[Math.floor(Math.random() * Math.floor(quotes.length))];
      // console.log(q);
      setQuote(q);
    }, 90 * 100);
    return () => {
      clearInterval(i_id);
    };
  }, []);

  return (
    <>
      <p style={{ margin: "1rem" }}>
        <span role="img" aria-label="Light Bulb">
          💡
        </span>
        {quote}
      </p>
    </>
  );
}

export default Tips;
