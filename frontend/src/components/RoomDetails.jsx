import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { hotelService, reservationService } from '../services/api';
import '../estilos/RoomDetails.css';

const RoomDetails = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [room, setRoom] = useState(null);
  const [loading, setLoading] = useState(true);
  const [dates, setDates] = useState({
    checkIn: '',
    checkOut: ''
  });

  useEffect(() => {
    fetchRoomDetails();
  }, [id]);

  const fetchRoomDetails = async () => {
    try {
      const data = await hotelService.getById(id);
      setRoom(data);
    } catch (err) {
      console.error('Error al cargar detalles:', err);
    }
    setLoading(false);
  };

  const handleReservation = async () => {
    const token = localStorage.getItem('token');
    
    if (!token) {
      navigate('/login');
      return;
    }

    try {
      await reservationService.create({
        roomId: id,
        checkIn: dates.checkIn,
        checkOut: dates.checkOut
      });
      navigate('/congrats');
    } catch (err) {
      console.error('Error:', err);
      alert('Error al crear la reserva');
    }
  };

  if (loading) return <p className="loading">Cargando...</p>;
  if (!room) return <p className="no-results">Hotel no encontrado</p>;

  return (
    <div className="room-details-container">
      <button onClick={() => navigate('/')} className="btn-back">
        ← Volver
      </button>
      
      <img 
        src={room.image || '/placeholder-room.jpg'} 
        alt={room.name}
        className="room-image-main"
      />
      
      <h1>{room.name}</h1>
      <p className="room-description">{room.description}</p>
      
      <div className="room-info">
        <p><strong>Capacidad:</strong> {room.capacity} personas</p>
        <p className="room-price-large"><strong>Precio:</strong> ${room.price} por noche</p>
        <p><strong>Servicios:</strong> {room.amenities?.join(', ') || 'WiFi, TV, Aire acondicionado'}</p>
      </div>

      <div className="reservation-box">
        <h3>Reservar esta habitación</h3>
        <div className="date-group">
          <label>Check-in:</label>
          <input 
            type="date"
            value={dates.checkIn}
            onChange={(e) => setDates({...dates, checkIn: e.target.value})}
          />
        </div>
        <div className="date-group">
          <label>Check-out:</label>
          <input 
            type="date"
            value={dates.checkOut}
            onChange={(e) => setDates({...dates, checkOut: e.target.value})}
          />
        </div>
        <button 
          onClick={handleReservation}
          disabled={!dates.checkIn || !dates.checkOut}
          className="btn-reserve"
        >
          Confirmar Reserva
        </button>
      </div>
    </div>
  );
};

export default RoomDetails;