import type { Quiz } from "../model/quiz";

const API_BASE_URL = process.env.NODE_ENV === 'production' 
  ? 'https://your-vercel-deployment-url.vercel.app/api' 
  : 'http://localhost:3000/api';

export class ApiService {
    async getQuizById(id: string): Promise<Quiz | null> {
        let response = await fetch(`${API_BASE_URL}/quizzes/${id}`);
        if (!response.ok) {
            return null;
        }

        let json = await response.json();
        return json;
    }

    async getQuizzes(): Promise<Quiz[]> {
        let response = await fetch("/api/quizzes");
        if (!response.ok) {
            alert("Failed to fetch quizzes!");
            return [];
        }

        let json = await response.json();
        return json;
    }

    async saveQuiz(quizId: string, quiz: Quiz) {
        let response = await fetch(`/api/quizzes/${quizId}`, {
            method: "PUT",
            body: JSON.stringify(quiz),
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (!response.ok) {
            alert("Failed to save quiz!");
            return;
        }
    }
}

export const apiService = new ApiService();
