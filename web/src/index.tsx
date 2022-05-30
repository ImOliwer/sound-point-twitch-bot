import React from 'react';
import ReactDOM from 'react-dom/client';
import './style/global.css'
import Deployments from './Deployments';
import reportWebVitals from './reportWebVitals';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import NotFound from './NotFound';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <BrowserRouter>
    <Routes>
      <Route path="/sound-deployments" element={<Deployments />} />
      <Route path="/*" element={<NotFound/>} />
    </Routes>
  </BrowserRouter>
);

reportWebVitals();
