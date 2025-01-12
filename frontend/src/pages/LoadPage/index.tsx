import { signal } from "@preact/signals";
import { useRoute } from "preact-iso";

export function LoadPage() {
	const route = useRoute();
	const id = route.path.slice(1);

	const cpu = signal();
	const mem = signal();
	const uptime = signal();

	return (
		<table>
			<thead>
				<tr>
					<th>CPU</th>
					<th>MEM</th>
					<th>Uptime</th>
				</tr>
			</thead>

			<tbody>
				<tr>
					<td>{cpu}</td>
					<td>{mem}</td>
					<td>{uptime}</td>
				</tr>
			</tbody>
		</table>
	);
}
