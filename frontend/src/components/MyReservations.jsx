import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import ReservationCard from './ReservationCard';
import { reservationService } from '../services/api';
import '../estilos/MyReservations.css';

const MyReservations = () => {
  const [reservations, setReservations] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    fetchReservations();
  }, []);

  const fetchReservations = async () => {
    const token = localStorage.getItem('token');
    
    if (!token) {
      navigate('/login');
      return;
    }

    try {
      const data = await reservationService.getMyReservations();
      setReservations(data);
    } catch (err) {
      console.error('Error al cargar reservas:', err);
    }
    setLoading(false);
  };

  const handleCancel = async (id) => {
    if (!confirm('¿Seguro que quieres cancelar esta reserva?')) return;

    try {
      await reservationService.delete(id);
      setReservations(reservations.filter(r => r.id !== id));
    } catch (err) {
      console.error('Error:', err);
      alert('Error al cancelar la reserva');
    }
  };

  return (
    <div className="my-reservations-container">
      <h2>Mis Reservas</h2>
      {loading ? (
        <p className="loading">Cargando...</p>
      ) : (
        <>
          {reservations.length > 0 ? (
            reservations.map(reservation => (
              <ReservationCard 
                key={reservation.id} 
                reservation={reservation}
                onCancel={handleCancel}
              />
            ))
          ) : (
            <div className="no-reservations">
              <p>No tienes reservas activas</p>
              <a href="/" className="btn-browse">Explorar habitaciones</a>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default MyReservations;