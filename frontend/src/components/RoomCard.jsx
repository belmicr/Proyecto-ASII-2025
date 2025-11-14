import { useNavigate } from 'react-router-dom';
import '../estilos/RoomCard.css';

const RoomCard = ({ room }) => {
  const navigate = useNavigate();

  return (
    <div className="room-card">
      <img 
        src={room.image || '/placeholder-room.jpg'} 
        alt={room.name}
        className="room-card-image"
      />
      <div className="room-card-content">
        <h3>{room.name}</h3>
        <p>{room.description}</p>
        <p><strong>Capacidad:</strong> {room.capacity} personas</p>
        <p className="room-price">${room.price} por noche</p>
        <button 
          onClick={() => navigate(`/room/${room.id}`)}
          className="btn-details"
        >
          Ver Detalles
        </button>
      </div>
    </div>
  );
};

export default RoomCard;