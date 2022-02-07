import React, { useEffect, useState } from "react";
import "./App.css";
import Tips from "./components/Tips";
import logo from "./cvoid19logo.png";
import { db } from "./service/firebaseService";
import style from "./style.module.css";

function App() {
  return (
    <>
      <div className={style.page_container}>
        <header className={`${style.content_parent}, ${style.blue}`}>
          <div className={style.box_header}>
            <img src={logo} alt="logo" width="100" height="100" />
          </div>
        </header>
        <Crowd />
        <Tables />

        <footer className={style.purple}>
          <Tips />
        </footer>
      </div>
    </>
  );
}

function Crowd(props) {
  const [enter, setEnter] = useState(0);
  const [close, setClose] = useState(0);
  const [crowd, setCrowd] = useState(0);

  useEffect(() => {
    db.ref("cvoid-bar/door-enter").on("value", function (snapshot) {
      setEnter(snapshot.numChildren());
    });
    db.ref("cvoid-bar/door-exit").on("value", function (snapshot) {
      setClose(snapshot.numChildren());
    });
    setCrowd(enter - close);
  }, [enter, close]);

  return (
    <>
      <main className={`${style.content_parent}, ${style.coral}`}>
        <div className={style.box}>
          <span>PEOPLE INSIDE</span>
          <span className={`${style.padded} , ${style.green}`}>{crowd}</span>
        </div>
      </main>
    </>
  );
}

function Tables(props) {
  return (
    <>
      <main className={`${style.table_parent} , ${style.coral}`}>
        <Table arrayOfValue={[0.5, 0.8, 0.6, 0.7]} number={1} />
        <Table arrayOfValue={[0.5, 0.2, 0.6, 0.7]} number={2} />
        <Table arrayOfValue={[0.5, 0.2, 0.6, 0.7]} number={3} />
        <Table arrayOfValue={[0.5, 0.2, 0.6, 0.7]} number={4} />
      </main>
    </>
  );
}

function Table(props) {
  // const [tabValue, setTabValue] = useState(props.arrayOfValue);
  const [tabValue, setTabValue] = useState([0, 0, 0, 0]);

  useEffect(() => {
    const timer = setTimeout(() => {
      db.ref(`cvoid-bar/tables/table_${props.number}`).once(
        "value",
        function (snapshot) {
          const snap = Object.values(snapshot.val());
          setTabValue(
            snap.map((item, index) => {
              return item.value;
            })
          );
        }
      );
    }, 20 * 1000);
    return () => clearTimeout(timer);
  });

  const color = {
    green: "#A1D884",
    red: "#F8485E",
  };

  const allSeatOccupated = `${color.red},${color.red},${color.red},${color.red}`;

  let seatOccupated = tabValue.map((item, index) => {
    if (item > 0.25) {
      return color.red;
    }
    return color.green;
  });

  let tableStyle = {
    border: "7px solid",
    borderColor: seatOccupated.join(" "),
  };

  return (
    <>
      {allSeatOccupated.toString() === seatOccupated.toString() ? (
        <div className={`${style.padded} , ${style.pink}`} style={tableStyle}>
          {props.number}
        </div>
      ) : (
        <div className={`${style.padded} , ${style.green}`} style={tableStyle}>
          {props.number}
        </div>
      )}
    </>
  );
}
export default App;
