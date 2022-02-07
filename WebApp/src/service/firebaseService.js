import { firebaseConfig } from "../config/FirebaseConfig";
import firebase from "firebase/app";
import "firebase/auth";
import "firebase/database";

firebase.initializeApp(firebaseConfig);

export const auth = firebase.auth;
export const db = firebase.database();
