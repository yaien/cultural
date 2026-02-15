document.addEventListener("alpine:init", () => {
  Alpine.data("roles", ({ currentUserId }) => ({
    currentUserId,
    roles: [],
    loading: true,
    deleting: true,
    selected: null,
    modals: {
      invitation: false,
      edition: false,
      deletion: false,
    },
    init() {
      this.fetch();
    },
    async fetch() {
      try {
        this.loading = true;
        const res = await fetch("/dashboard/api/roles");
        this.roles = (await res.json()) || [];
      } catch (error) {
        console.error(error);
      } finally {
        this.loading = false;
      }
    },

    async openDeleteModal(role) {
      this.selected = role;
      this.modals.deletion = true;
    },

    onDeleted() {
      this.modals.deletion = false;
      this.fetch();
    },
  }));

  Alpine.data("invitation", () => ({
    loading: false,
    userDisplayName: "",
    userEmail: "",

    async invite() {
      try {
        this.loading = true;
        const res = await fetch("/dashboard/api/invitations", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            groupId: null,
            permissions: ["*"],
            name: "Admin",
            userDisplayName: this.userDisplayName,
            userEmail: this.userEmail,
          }),
        });

        if (!res.ok) {
          const data = await res.json();
          throw Error(data.error);
        }

        this.$dispatch("toast", { message: "Invitación enviada correctamente a el correo " + this.userEmail });
        this.$dispatch("submitted");
      } catch (error) {
        switch (error.message) {
          case "user_already_exist":
            this.$dispatch("toast", { message: "El correo ya pertenece a un rol asignado", type: "warning" });
            break;
          default:
            this.$dispatch("toast", { message: "Error al enviar la invitación", type: "danger" });
        }
      }
    },
  }));

  Alpine.data("deletion", ({ selected }) => ({
    loading: false,
    selected: selected,

    async confirmDelete() {
      try {
        this.loading = true;
        const res = await fetch(`/dashboard/api/roles/${this.selected.id}`, {
          method: "DELETE",
        });

        if (!res.ok) {
          const data = await res.json();
          throw Error(data.error);
        }

        this.$dispatch("toast", { message: "Rol eliminado correctamente" });
        this.$dispatch("deleted");
      } catch (error) {
        this.$dispatch("toast", { message: "Error al eliminar el rol", type: "danger" });
      }
    },
  }));
});
