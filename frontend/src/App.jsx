import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './components/Login';
import Register from './components/Register';
import RoomList from './components/RoomList';
import RoomDetails from './components/RoomDetails';
import MyReservations from './components/MyReservations';

const Congrats = () => (
  <div style={{ textAlign: 'center', marginTop: '100px' }}>
    <h1 style={{ color: '#28a745' }}>¡Reserva Exitosa! ✓</h1>
    <p>Tu reserva ha sido confirmada correctamente</p>
    <a href="/my-reservations" style={{ 
      display: 'inline-block', 
      marginTop: '20px', 
      padding: '10px 20px', 
      background: '#cce206', 
      color: '#000', 
      textDecoration: 'none',
      borderRadius: '4px',
      fontWeight: '600'
    }}>
      Ver mis reservas
    </a>
  </div>
);

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/" element={<RoomList />} />
        <Route path="/room/:id" element={<RoomDetails />} />
        <Route path="/my-reservations" element={<MyReservations />} />
        <Route path="/congrats" element={<Congrats />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
}

export default App;