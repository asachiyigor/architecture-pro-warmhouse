"""Database configuration and session management"""

from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker
from app.core.config import settings
from app.models.telemetry import Base

# Create async engine
engine = create_async_engine(
    settings.DATABASE_URL,
    echo=settings.DEBUG,
    future=True
)

# Create async session factory
async_session_factory = sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False
)


async def init_db():
    """Initialize database - create tables"""
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)


async def get_db_session() -> AsyncSession:
    """Dependency for getting database session"""
    async with async_session_factory() as session:
        try:
            yield session
        finally:
            await session.close()