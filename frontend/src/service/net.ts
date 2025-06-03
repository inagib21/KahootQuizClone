import type { Player, QuizQuestion } from "../model/quiz";

export enum PacketTypes {
    Connect,
    HostGame,
    QuestionShow,
    ChangeGameState,
    PlayerJoin,
    StartGame,
    Tick,
    Answer,
    PlayerReveal,
    Leaderboard,
    PlayerDisconnect
}

export enum GameState {
    Lobby,
    Play,
    Intermission,
    Reveal,
    End
}

export interface Packet {
    id: PacketTypes;
}

export interface HostGamePacket extends Packet {
    quizId: string;
}

export interface ChangeGameStatePacket extends Packet {
    state: GameState;
}

export interface PlayerJoinPacket extends Packet {
    player: Player;
}

export interface TickPacket extends Packet {
    tick: number;
}

export interface PlayerDisconnectPacket extends Packet {
    playerId: string;
}

export interface ConnectPacket extends Packet {
    code: string;
    name: string;
}

export interface QuestionShowPacket extends Packet {
    question: QuizQuestion;
}

export interface QuestionAnswerPacket extends Packet {
    question: number;
}

export interface PlayerRevealPacket extends Packet {
    points: number;
}

export interface LeaderboardEntry {
    name: string;
    points: number;
}

export interface LeaderboardPacket extends Packet {
    points: LeaderboardEntry[];
}

export class NetService {

    private eventSource!: EventSource;

    private onPacketCallback?: (packet: any) => void;

    connect(){
        this.eventSource = new EventSource("/api/events");
        
        this.eventSource.onopen = () => {
            console.log("SSE connection opened");
        };

        this.eventSource.onerror = (err) => {
            console.error("EventSource failed:", err);
        };

        this.eventSource.onmessage = (event: MessageEvent) => {
            console.log("SSE message received:", event.data);
            try {
                const packetData = JSON.parse(event.data);
                if (packetData && typeof packetData.id !== 'undefined') {
                    if(this.onPacketCallback) {
                        this.onPacketCallback(packetData);
                    }
                } else {
                    if (packetData.type === "connected") {
                        console.log("SSE Server Connection Message:", packetData.message);
                    }
                }
            } catch (e) {
                console.error("Error parsing SSE message data:", e);
            }
        };
    }

    onPacket(callback: (packet: Packet) => void){
        this.onPacketCallback = callback;
    }

    sendPacket(packet: Packet) {
		console.warn("sendPacket (SSE): This action needs to be refactored to use HTTP requests.", packet);
	}

}
