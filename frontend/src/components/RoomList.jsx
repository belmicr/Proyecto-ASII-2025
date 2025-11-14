import { useState, useEffect } from 'react';
import SearchBar from './SearchBar';
import RoomCard from "./RoomCard";
import { searchService } from '../services/api';
import '../estilos/RoomList.css';

const RoomList = () => {
  const [rooms, setRooms] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);

  const fetchRooms = async (searchQuery = '') => {
    setLoading(true);
    try {
      const data = await searchService.search(searchQuery, page);
      setRooms(data.results || data.rooms || data);
    } catch (err) {
      console.error('Error al cargar hoteles:', err);
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchRooms();
  }, [page]);

  const handleSearch = (query) => {
    setPage(1);
    fetchRooms(query);
  };

  return (
    <div>
      <SearchBar onSearch={handleSearch} />
      <div className="room-list-container">
        <h2>Hoteles Disponibles</h2>
        {loading ? (
          <p className="loading">Cargando...</p>
        ) : (
          <>
            {rooms.length > 0 ? (
              rooms.map(room => <RoomCard key={room.id} room={room} />)
            ) : (
              <div className="no-results">
                <p>No se encontraron hoteles</p>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default RoomList;