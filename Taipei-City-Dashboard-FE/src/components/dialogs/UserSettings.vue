<!-- Developed by Taipei Urban Intelligence Center 2023-2024-->

<script setup>
import { storeToRefs } from "pinia";
import { useDialogStore } from "../../store/dialogStore";
import { useAuthStore } from "../../store/authStore";

import DialogContainer from "./DialogContainer.vue";

const dialogStore = useDialogStore();
const authStore = useAuthStore();

const { editUser } = storeToRefs(authStore);

function handleClose() {
	dialogStore.hideAllDialogs();
	authStore.editUser = authStore.user;
}

function parseTime(time) {
	return time.slice(0, 19).replace("T", " ");
}
</script>

<template>
	<DialogContainer @onClose="handleClose" :dialog="`userSettings`">
		<div class="usersettings">
			<h2>用戶設定</h2>
			<label> 用戶名稱 </label>
			<input
				v-model="editUser.name"
				:minlength="1"
				:maxlength="10"
				required
			/>
			<label> 用戶帳號 </label>
			<input :value="editUser.account" :minlength="1" disabled />
			<label> 用戶類型 </label>
			<input
				:value="editUser.isAdmin ? '管理員' : '一般用戶'"
				disabled="true"
				required
			/>
			<label> 最近登入時間 </label>
			<input :value="parseTime(editUser.login_at)" disabled />
		</div>
	</DialogContainer>
</template>

<style scoped lang="scss">
.usersettings {
	width: 300px;
	display: flex;
	flex-direction: column;

	label {
		margin: 8px 0 4px;
		font-size: var(--font-s);
		color: var(--color-complement-text);
	}
}
</style>
