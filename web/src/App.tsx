import { Routes, Route, Navigate } from 'react-router-dom';
import AttemptHistory from './pages/AttemptHistory';
import AttemptDetails from './pages/AttemptDetails';
import AuthLayout from './layouts/AuthLayout';
import DashboardLayout from './layouts/DashboardLayout';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import DictationList from './pages/DictationList';
import CreateDictation from './pages/CreateDictation';
import PlayDictation from './pages/PlayDictation';
import './App.css';

function App() {
  return (
    <Routes>
      <Route path="/auth" element={<AuthLayout />}>
        <Route path="login" element={<Login />} />
        <Route path="register" element={<Register />} />
      </Route>

      <Route path="/" element={<DashboardLayout />}>
        <Route index element={<Dashboard />} />
        <Route path="dictations" element={<DictationList />} />
        <Route path="dictations/new" element={<CreateDictation />} />
        <Route path="play/:id" element={<PlayDictation />} />
        <Route path="attempts" element={<AttemptHistory />} />
        <Route path="attempts/:id" element={<AttemptDetails />} />
      </Route>

      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

export default App;
