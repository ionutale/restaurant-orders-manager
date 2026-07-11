export interface FloorPlanTable {
	id: number;
	name: string;
	capacity: number;
	x: number;
	y: number;
	label: string | null;
	status?: 'free' | 'occupied';
	group_id?: number;
	group_name?: string;
	party_size?: number;
}
